

### 
Simplify  process build and push your docker image with any tools CI/CD like Gitlab-ci, Jenkins etc



### Usage 
```bash
imam@imam-mv:~/learn-golang/docker-api$ go build -o go-build-image
imam@imam-mv:~/learn-golang/docker-api$ ./go-build-image 
{"stream":"Step 1/5 : FROM alpine:latest"}
{"stream":"\n"}
{"stream":" ---\u003e d4ff818577bc\n"}
{"stream":"Step 2/5 : WORKDIR /app"}
{"stream":"\n"}
{"stream":" ---\u003e Using cache\n"}
{"stream":" ---\u003e cebd6b2e2c6b\n"}
{"stream":"Step 3/5 : RUN echo \"hello\" \u003e bandung.txt"}
{"stream":"\n"}
{"stream":" ---\u003e Running in 026abaf95ab8\n"}
{"stream":"Removing intermediate container 026abaf95ab8\n"}
{"stream":" ---\u003e 6b60236f0c68\n"}
{"stream":"Step 4/5 : RUN echo \"bandung\" \u003e jakarta.txt"}
{"stream":"\n"}
{"stream":" ---\u003e Running in 0f08ec39130d\n"}
{"stream":"Removing intermediate container 0f08ec39130d\n"}
{"stream":" ---\u003e 7de7761fb55a\n"}
{"stream":"Step 5/5 : RUN echo \"semarang\" \u003e semarang.txt"}
{"stream":"\n"}
{"stream":" ---\u003e Running in 722169c7ecce\n"}
{"stream":"Removing intermediate container 722169c7ecce\n"}
{"stream":" ---\u003e c68a334b31fb\n"}
{"stream":"Successfully built c68a334b31fb\n"}
{"stream":"Successfully tagged registry.xxx.id/example/jakarta:v12.0.0\n"}
2021/07/04 16:37:55 Has successfuly build registry.xxx.id/example/jakarta:v12.0.0 
The push refers to repository [registry.xxx.id/example/jakarta]
c617b6bf1a7a: Preparing
6fdeec4f4cc1: Preparing
0b5af8654dcf: Preparing
9b558a5aaf27: Preparing
72e830a4dff5: Preparing
9b558a5aaf27: Layer already exists
72e830a4dff5: Layer already exists
6fdeec4f4cc1: Pushed
0b5af8654dcf: Pushed
c617b6bf1a7a: Pushed
v12.0.0: digest: sha256:b65cb0a17f021714b4ec1374fbd5a2b2db9388dfa1a1d196e72dfc5162f9b977 size: 1356
The push refers to repository [registry.xxx.id/example/jakarta]
c617b6bf1a7a: Preparing
6fdeec4f4cc1: Preparing
0b5af8654dcf: Preparing
9b558a5aaf27: Preparing
72e830a4dff5: Preparing
9b558a5aaf27: Layer already exists
0b5af8654dcf: Layer already exists
72e830a4dff5: Layer already exists
6fdeec4f4cc1: Layer already exists
c617b6bf1a7a: Layer already exists
v12.0.0: digest: sha256:b65cb0a17f021714b4ec1374fbd5a2b2db9388dfa1a1d196e72dfc5162f9b977 size: 1356
The push refers to repository [registry.xxx.id/example/jakarta]
c617b6bf1a7a: Preparing
6fdeec4f4cc1: Preparing
0b5af8654dcf: Preparing
9b558a5aaf27: Preparing
72e830a4dff5: Preparing
6fdeec4f4cc1: Layer already exists
0b5af8654dcf: Layer already exists
72e830a4dff5: Layer already exists
9b558a5aaf27: Layer already exists
c617b6bf1a7a: Layer already exists
v12.0.0: digest: sha256:b65cb0a17f021714b4ec1374fbd5a2b2db9388dfa1a1d196e72dfc5162f9b977 size: 1356
The push refers to repository [registry.xxx.id/example/jakarta]
c617b6bf1a7a: Preparing
6fdeec4f4cc1: Preparing
0b5af8654dcf: Preparing
9b558a5aaf27: Preparing
72e830a4dff5: Preparing
c617b6bf1a7a: Layer already exists
9b558a5aaf27: Layer already exists
6fdeec4f4cc1: Layer already exists
0b5af8654dcf: Layer already exists
72e830a4dff5: Layer already exists
v12.0.0: digest: sha256:b65cb0a17f021714b4ec1374fbd5a2b2db9388dfa1a1d196e72dfc5162f9b977 size: 1356
2021/07/04 16:37:57 Has Successfully push to registry 
2021/07/04 16:37:57 Clean image successfully

```


### Reference ###
- https://gowalker.org/github.com/fsouza/go-dockerclient


