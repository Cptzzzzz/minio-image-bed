# minio-image-bed

> 使用`minio`搭建`typora`图床

## 搭建`minio`服务

### Docker运行

首先需要搭建`minio`对象存储服务，需要通过环境变量设置`root`用户的账号密码

```
docker run -d --name minio-server -p 9000:9000 -p 9001:9001 --env MINIO_ROOT_USER="minio-root-user" --env MINIO_ROOT_PASSWORD="minio-root-password"  bitnami/minio:latest
```

随后进入 `IP:9001`，使用上述设置的账号密码登录进入`minio dashboard`

### 创建用户

单击User，新建一个访问用户，随后为其添加全部权限

![1681181851](http://storage.cptz.space/typora/1681205096.png)

随后生成该用户的`Service Account`

![1681181953](http://storage.cptz.space/typora/1681205122.png)

![1681181926](http://storage.cptz.space/typora/1681205150.png)

记住生成的账号密码

### 设置存储`bucket`

随后点击左侧`bucket`，新建一个`bucket`，取名随意，例如下图创建了名为`images`和`typora`的两个`bucket`

![1681181383](http://storage.cptz.space/typora/1681205178.png)

随后选择你想存储的`bucket`，点击`manage`

在`summary`中设置`Access Policy`为`public`

在`Access Rules`中添加一条

![1681182097](http://storage.cptz.space/typora/1681205197.png)

至此，我们的`minio`服务已经配置完成

访问`minio`中文件的路径为 `BaseUrl/BucketName/filename`

例如访问本地`9000`端口`bucket`中的`filename`文件

URL:`127.0.0.1：9000/bucket/filename`

## 服务端存储图片

> 这一部分的主要工作是接收上传的图片，使用上述的`Service Account`连接`minio`，并将图片存入

参考代码

`server.go`

```go
package main

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"net/http"
	"path"
	"time"
)

var minioClient *minio.Client

const BucketName = "typora"
const Domain = "127.0.0.1:9000"

func connect() {
	minioClient = nil
	endpoint := Domain
	accessKeyID := "MQTYX33BJR24UQKWBVKY"
	secretAccessKey := "elQ+3ngJ9+6C+XlVhu1FPAZX3PRMyFHEinmC3Jo7"
	useSSL := false
	minio, _ := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	minioClient = minio
}

func main() {
	connect()
	http.HandleFunc("/upload", upload)
	http.ListenAndServe("0.0.0.0:8080", nil)
}

func generateFileName(name string) string {
	return fmt.Sprintf("%v%s", time.Now().Unix(), name)
}

func upload(w http.ResponseWriter, r *http.Request) {
	file, fileHeader, _ := r.FormFile("file")
	filename := generateFileName(path.Ext(fileHeader.Filename))
	minioClient.PutObject(context.Background(), BucketName, filename,
		file, fileHeader.Size, minio.PutObjectOptions{})
	w.Write([]byte(fmt.Sprintf("%s/%s/%s", Domain, BucketName, filename)))
}

```

这份代码监听`8080`端口，同时接收`Http`请求中传入的文件，将文件名按时间生成后存入`minio`，并把`URL`返回给客户端

打包镜像

```shell
docker build -t minio-image-server:1 .
```

填写信息，运行镜像

```shell
docker run -d -p 8080:8080 -e Endpoint="" -e AccessKey="" -e SecretKey="" -e Bucket="" minio-image-server:1
```

## 客户端上传图片

> 这一部分的主要功能是对接`Typora`，封装`Http`请求完成上传图片

根据`Typora`中自定义上传图片的逻辑

![image-20230411111255470](http://storage.cptz.space/typora/1681182847.png)

在`Custom Command`中填入所需执行的命令，`Typora`会在命令后补全文件绝对路径并调用，因此需要从命令行参数中读取文件名，随后封装`Http`请求，最后把图片`URL`输出到标准输出，`Typora`会根据标准输出的结果替换图片路径。

参考代码

`client.go`

```go
package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

var UploadUrl = "localhost:8080/api/upload"

func main() {
	content, _ := os.ReadFile(os.Args[1])
	buf := new(bytes.Buffer)
	bodyWriter := multipart.NewWriter(buf)
	w, _ := bodyWriter.CreateFormFile("file", os.Args[1])
	w.Write(content)
	bodyWriter.Close()
	resp, err := http.Post(UploadUrl, bodyWriter.FormDataContentType(), buf)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	fmt.Println(string(data))
}

```

最后使用`go build -o upload.exe client.go`，将可执行文件绝对目录填入`Typora`中即可

## 补充

如果你有域名的话，可以使用`nginx`来为你的`minio`图床设置域名。

实际运行时，需要把代码中的`URL`换成真实的`URL`
