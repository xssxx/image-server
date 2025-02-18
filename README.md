# image-server

## Test
Upload image
```bash
curl -X POST localhost:8000/upload -F "image=path_image.jpg"  
```
Get image
```bash
curl localhost:8000/image/path_image.jpg
```
