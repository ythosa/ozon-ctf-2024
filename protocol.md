Структура пакета:
    offset: +0,         size: 4,            data: "PDR\x00",    desc: protocol magic
    offset: +4,         size: 4,            data: <bytes>,      desc: idx of block
    offset: +8,         size: 4,            data: <bytes>,      desc: size of block
    offset: +12,        size: <optional>,   data: <bytes>,      desc: block data
    offset: +optional,  size: 4,            data: <bytes>,      desc: crc32 of block data

Специальные опции:
    пиксели изображения ключа передаются на чётный порт (например, 1024, 4040, 1234)
    пиксели шифрованного изображение передаются на нечётный порт (1011, 3035, 8485)

Схема шифрования:
Входное изображение шифруется с помощью изображения ключа. 
Каждые 48 байт (16 пикселей) входного изображения шифруются с помощью алгоритма AES-ECB на ключе из 16 байт.
16 байт для ключа получаются путём взятия первых байт из блока 48 байт ключевого изображения. 
Изображения равны по разрешению и содержат одинаковое количество пикселей.
Размер изображений 1280х720 пикселей. 
Ниже приведён код шифрования изображения, а также код перевода двумерного массива пикселей в 
линейный массив который передаётся по сети.

```Python
def get_pixels_as_linear_list(key_img_name: str):
    src_img = Image.open(key_img_name)
    src_pixels = src_img.load()

    w, h = src_img.size
    data = []
    
    for i in range(w):
        for j in range(h):
            data += [src_pixels[i, j][0], src_pixels[i, j][1], src_pixels[i, j][2]]

    src_img.close()
    return bytes(data)

 def main(args: list):
    if len(sys.argv) > 2:
        src_file = sys.argv[1]
        key_file = sys.argv[2]
    else:
        print("[-] Usage: ./gen.py <src-file> <key-file>")
        sys.exit(-1)

    img = Image.open(src_file)
    print(img.mode)

    src_pixels = get_pixels_as_linear_list(src_file)
    key_pixels = get_pixels_as_linear_list(key_file)

    res_pixels = b''

    for i in range(0, len(src_pixels), 48):
        ctx = AES.new(key_pixels[i:i+16], AES.MODE_ECB)
        res_pixels += ctx.encrypt(src_pixels[i:i+48])

    res_img = Image.new("RGBA", (1280, 720))
    off = 0
    for i in range(res_img.size[0]):
        for j in range(res_img.size[1]):
            pixel = (res_pixels[off], res_pixels[off + 1], res_pixels[off + 2])
            off += 3
            res_img.putpixel((i,j), pixel)

    res_img.save("res.png")
    res_img.close()
```
