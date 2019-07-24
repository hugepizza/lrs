# 导入模块
from wxpy import *
import time
# 初始化机器人，扫码登陆
bot = Bot()

def task():
	my_friend = bot.friends().search('💙',sex=FEMALE)[0]
	my_friend.send('狼人杀关键字监控')
	my_friend.send_image('pc_360.png')
	my_friend.send_image('pc_baidu.png')
	my_friend.send_image('pc_sougou.png')
	my_friend.send_image('360.png')
	my_friend.send_image('baidu.png')
	my_friend.send_image('sougou.png')
	my_friend.send_image('shenma.png')
	my_friend.send_image('360.png')

time.sleep(120)

while True:    
    print("wechat sending")    
    task()      
    time.sleep(3600*7.5)

