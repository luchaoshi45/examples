# webui 🎇🎇🎇

## 1 Enter the directory
```shell
cd sound-equipment-fault-detection/webui
```

## 2 Use the config file in the cloud to replace the config in this folder
```shell
cp ~/.kube/config ./
```

## 3 Build the image
```shell
docker build -f Dockerfile -t webui .
docker images
```
## 4 Deploy webui ✅
```shell 
kubectl apply -f resource/deployment.yaml
```