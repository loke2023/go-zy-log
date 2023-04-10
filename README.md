
##环境要求
首次安装，必须安装docker, docker-compose, ffmpeg, go, git, python3.7.4, openvc-python


##主要步骤
先启动redis： docker-compose up -d

先编译，go build ./main.go 

启动：
   编译完之后启动，执行 ./main & 
   退出时，输入exit， 如果直接关控制台，可能会导致后台服务进程退出


更新：
   用 ps -ef|grep main 看下进程ID，然后kill掉
   然后再进入/root/go/src/go-zy-log目录，执行 ./main & 
   退出时，输入exit， 如果直接关控制台，可能会导致后台服务进程退出


停止
   用 ps -ef|grep main 看下进程ID，然后kill掉



##安装docker

安装docker-ce
1、切换为root用户

2、卸载旧版本（如果安装过旧版本的话）
$ yum remove -y docker \
  docker-client \
  docker-client-latest \
  docker-common \
  docker-latest \
  docker-latest-logrotate \
  docker-logrotate \
  docker-selinux \
  docker-engine-selinux \
  docker-engine
  
3、安装需要的软件包
 yum-util 提供 yum-config-manager 功能
 另外两个是 devicemapper 驱动依赖的
$ yum install -y yum-utils \ 
device-mapper-persistent-data \ 
lvm2 

4、设置yum源
由于官网的源太慢了，这里可以使用阿里云Docker Yum源替代： 
$ sudo yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
$ sudo yum makecache fast 

5、安装最新版本docker
$ yum install -y docker-ce

6、启动docker
$ systemctl start docker.service

7、设置开机自启动
$ systemctl enable docker.service

8、验证安装是否成功（有client和service两部分表示docker安装启动都成功了）
$ docker version 

###安装docker-compose
github的地址下载太慢了，国内可以使用http://get.daocloud.io/#install-compose网站上面的地址。

首先下载docker-compose：

curl -L https://get.daocloud.io/docker/compose/releases/download/1.25.4/docker-compose-`uname -s`-`uname -m` > /usr/local/bin/docker-compose
然后赋权限：

chmod +x /usr/local/bin/docker-compose
最后查看版本：

docker-compose -v 
  
  
##安装ffmpeg
1.升级系统(也可以跳过这一步)
sudo yum install epel-release -y
sudo yum update -y

2.由于CentOS没有官方FFmpeg rpm软件包。但是，我们可以使用第三方YUM源（Nux Dextop）完成此工作。
sudo rpm --import http://li.nux.ro/download/nux/RPM-GPG-KEY-nux.ro
sudo rpm -Uvh http://li.nux.ro/download/nux/dextop/el7/x86_64/nux-dextop-release-0-5.el7.nux.noarch.rpm

3.安装FFmpeg 和 FFmpeg开发包
sudo yum install ffmpeg ffmpeg-devel -y

4.查看是否安装成功
ffmpeg -h

# 安装git
yum install -y git


##CentOS 安装 python3.7.4步骤：
sudo yum -y groupinstall "Development tools"
sudo yum -y install zlib-devel bzip2-devel openssl-devel ncurses-devel sqlite-devel readline-devel tk-devel gdbm-devel db4-devel libpcap-devel xz-devel libffi-devel

然后获取python3.7的安装包
wget https://www.python.org/ftp/python/3.7.4/Python-3.7.4.tar.xz
 
解压
tar -xvJf  Python-3.7.4.tar.xz

配置python3的安装目录并安装
cd Python-3.7.4
./configure --prefix=/usr/local/bin/python3
sudo make
sudo make install

创建软链接
ln -s /usr/local/bin/python3/bin/python3 /usr/bin/python3
ln -s /usr/local/bin/python3/bin/pip3 /usr/bin/pip3

验证是否成功
python3 -V
pip3 -V

安装numpy,此处采用pip安装
pip3 install numpy
****如果 报错: THESE PACKAGES DO NOT MATCH THE HASHES FROM THE REQUIREMENTS FILE.....
****改为pip3 install --upgrade numpy

安装opencv
pip3 install --upgrade opencv-python
****如果报错， 改为 pip3 install --upgrade opencv-python --no-use-pep517
****如果出现  No module named 'skbuild'
****更新 pip3 install --upgrade pip
****再次运行 pip3 install --upgrade opencv-python --no-use-pep517

测试是否安装成功
python3

进入python3环境

import cv2
没有反应，说明安装成功。



##安装go环境
yum install -y wget

下载
wget https://golang.google.cn/dl/go1.15.6.linux-amd64.tar.gz

解压
tar -C /usr/local -xzf go1.15.6.linux-amd64.tar.gz

将/usr/local/go/bin 目录添加至PATH环境变量
vim /etc/profile
export PATH=$PATH:/usr/local/go/bin

重新读取环境变量
source /etc/profile

查看版本
go version 

 