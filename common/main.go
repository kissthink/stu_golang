package main

import (
	"net"
	"net/http"
	"net/url"

	"strings"

	"net/http/cookiejar"
	"time"

	"context"
	"io/ioutil"
	"log"
	"os"
	//add by hyh
	//"golang.org/x/net/html"
	//"github.com/PuerkitoBio/goquery"
)

type HttpClient struct {
	ProxyAddr string
	Client    http.Client
}

const (
	DefaultIdleTimeout    = 75 * time.Second
	DefaultConnectTimeout = 45 * time.Second
)

// Conn wraps a net.Conn, and sets a deadline for every read
// and write operation.
type TimeoutConn struct {
	net.Conn
	IdleTimeout time.Duration
}

/////////////////////////
//test main
func main() {
	hc := NewHttpClient("https://p.xgj.me:27035")
	//hc := NewHttpClient("ip://192.168.1.102")
	//hc := NewHttpClient("ip://[2607:5300:60:6566::]")
	//url := "https://www.google.com"
	url := "https://play.google.com/store/search?q=facebook"
	req, e := http.NewRequest(
		"GET",
		url,
		nil,
	)
	resp, e := hc.Do(req)
	if e != nil {
		panic(e)
		return
	}
	defer resp.Body.Close()

	//1.获取到页面数据
	buf, e := ioutil.ReadAll(resp.Body)
	//print(string(buf))

	//2.保存查询页面到文件
	f, err1 := os.OpenFile("google_play_search.html", os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err1 != nil {
		panic(err1)
		return
	}
	defer f.Close()
	f.WriteString(string(buf))
	// Parse the HTML into nodes
	//root, e := html.Parse(res.Body)
	//if e != nil {
	//	return nil, e
	//}
	// Create and fill the document
	//goquery.NewDocumentFromResponse(resp)
}

/////////////////////////
//使用自定义出口协议,注意,前缀要全部使用小写
//如果是代理,那么使用 http:// 或者 https:// 类型的地址,如果使用出口 IP, 那么直接使用 ip:// 作为前缀
//如果使用ipv6, 那么使用`[]`把地址包起来
//例如:
//		http://14845132.xgj.me:27035
//		ip://192.168.1.12
//		ip://[2607:5300:60:6566::]

func MakeTransportX(addr string) (transport *http.Transport) {
	transport = new(http.Transport)
	transport.MaxIdleConnsPerHost = 16
	//disable verify ssl
	//transport.TLSClientConfig = &tls.Config{
	//	InsecureSkipVerify: true,
	//}
	var (
		localAddr string
		proxyUrl  string
	)
	if strings.HasPrefix(addr, "ip") {
		localAddr = addr[5:]
	} else if strings.HasPrefix(addr, "http") {
		proxyUrl = addr
	} else if addr != "" {
		log.Print("MakeTransportX, addr (", addr, ") have wrong format.")
	}
	transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		d := net.Dialer{Timeout: DefaultConnectTimeout}
		if localAddr != "" && localAddr[0] == '[' {
			//如果本地ip地址以"["开头, 那么是ipv6地址，强制使用 tcp6 拨号
			network = "tcp6"
		}
		lAddr, err := net.ResolveTCPAddr(network, localAddr+":0")
		if err != nil {
			return nil, err
		}
		d.LocalAddr = lAddr
		conn, err := d.DialContext(ctx, network, addr)
		if err != nil {
			return nil, err
		}
		return NewTimeoutConn(conn, DefaultIdleTimeout)
	}
	if proxyUrl != "" {
		u := url.URL{}
		urlProxy, e := u.Parse(proxyUrl)
		if e != nil {
			log.Print("set proxy failed, ", e)
		} else {
			transport.Proxy = http.ProxyURL(urlProxy)
		}
	}
	return transport

}

func NewTimeoutConn(conn net.Conn, idleTimeout time.Duration) (net.Conn, error) {
	c := &TimeoutConn{
		Conn:        conn,
		IdleTimeout: idleTimeout,
	}
	if c.IdleTimeout > 0 {
		deadline := time.Now().Add(idleTimeout)
		if e := c.Conn.SetDeadline(deadline); e != nil {
			return nil, e
		}
	}
	return c, nil
}
func (c *TimeoutConn) Read(b []byte) (int, error) {
	n, e := c.Conn.Read(b)
	if c.IdleTimeout > 0 && n > 0 && e == nil {
		err := c.Conn.SetDeadline(time.Now().Add(c.IdleTimeout))
		if err != nil {
			return 0, err
		}
	}
	return n, e
}

func (c *TimeoutConn) Write(b []byte) (int, error) {
	n, e := c.Conn.Write(b)
	if c.IdleTimeout > 0 && n > 0 && e == nil {
		err := c.Conn.SetDeadline(time.Now().Add(c.IdleTimeout))
		if err != nil {
			return 0, err
		}
	}
	return n, e
}

func NewHttpClient(proxyAddr string) *HttpClient {
	c := &HttpClient{ProxyAddr: proxyAddr}
	//Follow  时复制 http headerredirect
	c.Client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		for attr, val := range via[0].Header {
			if _, ok := req.Header[attr]; !ok {
				req.Header[attr] = val
			}
		}
		return nil
	}

	//	c.Client.Timeout = 75 * time.Second
	return c
}

func (hc *HttpClient) mkTransport() {
	if hc.Client.Transport != nil || hc.ProxyAddr == "" {
		return
	}
	hc.Client.Transport = MakeTransportX(hc.ProxyAddr)
}

func (hc *HttpClient) Do(req *http.Request) (resp *http.Response, err error) {
	hc.mkTransport()
	return hc.Client.Do(req)
}

func (hc *HttpClient) EnableCookie() {
	if hc.Client.Jar == nil {
		cookieJar, _ := cookiejar.New(nil)
		hc.Client.Jar = cookieJar
	}
}

func (hc *HttpClient) DisableCookie() {
	hc.Client.Jar = nil
}

func (hc *HttpClient) IsCookieEnabled() bool {
	return hc.Client.Jar != nil
}

func (hc *HttpClient) GetCookies(u *url.URL) []*http.Cookie {
	if hc.Client.Jar == nil {
		return nil
	}
	return hc.Client.Jar.Cookies(u)
}

func (hc *HttpClient) GetCookie(u *url.URL, key string) *http.Cookie {
	if hc.Client.Jar == nil {
		return nil
	}
	//	u, e := url.Parse(rawUrl)
	//	if  e != nil {
	//		return nil
	//	}
	for _, c := range hc.Client.Jar.Cookies(u) {
		if c.Name == "key" {
			return c
		}
	}
	return nil
}
