import requests
import pandas as pd
from bs4 import BeautifulSoup, Tag
import time


res = []
for i in range(53):
    url = 'http://www.xaprtc.com/queryContent_'+str(i)+'-jyxx.jspx?inDates=&ext=&origin=&channelId=271&projectCode=&title='
    print('crawling url ' + url)
    soup = BeautifulSoup(requests.get(url).content, 'html.parser')
    ul: Tag = soup.find('ul', class_="sx-lieul")
    for c in ul.children:
        if type(c) is not Tag:
            continue
        rs = c.find_all('p')
        dic = {
            'project': rs[0].get_text(),
            'date': rs[1].get_text()
        }
        res.append(dic)
    print('sleep 0.1 second')
    time.sleep(0.1)
df = pd.DataFrame(res)
print('save to csv file')
df.to_csv('data.csv')
