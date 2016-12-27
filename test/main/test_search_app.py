#coding:utf8
import urllib2
import json

def main():
    response = urllib2.urlopen('http://192.168.1.102:8000/search?q=wechat') 
    data = response.read()
    data = json.loads(data)
    # print data
    for app in data:
        print app
        print '-'*20

if __name__ == '__main__':
    main()