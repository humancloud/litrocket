# 插入单词表数据
import pymysql

# Connect To Mysql
conn = pymysql.connect(
    host='127.0.0.1',
    port=3306,
    user='root',
    password='linuxidc',
    database='litrocket')

# 创建游标
cursor = conn.cursor()

# 打开文件
v = open("v.dic")
w = open("w.dic")

# 创建字典
dict = {}

# 读取文件到字典中
for line in v.readlines():
    data = line[:-2] #去掉最后的两个字符,由于文件里面每一行后面有一个空格所以去掉最后两个
    list = data.split('#') #分割,返回分割后的列表
    dict[list[0]] = list[1]

for line in w.readlines():
    data = line[:-1] #这个文件里面每一行最后没有空格,只去掉换行就行
    list = data.split(' ') 
    dict[list[0]] = list[1]

# 关闭文件
v.close
w.close

# 根据字典数据创建SQL语句,并执行
for key in dict.keys():
    # Table dicts : {user_id,chinese,english}, user_id 表示软件自带单词,不是用户上传的个人单词
    sql = ("INSERT INTO dicts VALUES(%d,'%s','%s')") % (0,dict[key],key)
    cursor.execute(sql)

# 提交并关闭
conn.commit()
conn.close()