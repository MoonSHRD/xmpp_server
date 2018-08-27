package helpers

import (
    "image"
    "image/png"
    "image/color"
    "math/rand"
    "encoding/binary"
    "bytes"
    "encoding/base64"
)

const thumbSize = 128

func GenerateThumb(seed int64) string {
    //todo: fix this for calc funcs
    seed+=12742141
    
    buf := make([]byte, binary.MaxVarintLen64)
    n := binary.PutVarint(buf, seed)
    nameBytes := buf[:n]
    
    avatar := image.NewRGBA(image.Rect(0, 0, thumbSize, thumbSize))
    paintBG(avatar, calcBGColor(nameBytes))
    splatter(avatar, nameBytes, calcPixelColor(nameBytes))
    return encodeThumb(avatar)
}

//func SavePNG(avatar image.Image) {
//    file, err := os.Create(os.Args[1] + ".png")
//    err = png.Encode(file, avatar)
//
//    if err != nil {
//        panic(err)
//    }
//}

func encodeThumb(avatar image.Image) string {
    buf := new(bytes.Buffer)
    png.Encode(buf, avatar)
    send_s3 := buf.Bytes()
    imgBase64Str := base64.StdEncoding.EncodeToString(send_s3)
    return "data:image/png;base64,"+imgBase64Str
}

func splatter(avatar *image.RGBA, nameBytes []byte, pixelColor color.RGBA) {
    
    // A somewhat random number based on the username.
    var nameSum int64
    for i := range nameBytes {
        nameSum += int64(nameBytes[i])
    }
    
    // Use said number to keep random-ness deterministic for a given name
    rand.Seed(nameSum)
    
    // Make the "splatter"
    for y := 0; y < thumbSize; y++ {
        for x := 0; x < thumbSize; x++ {
            if ((x + y) % 2) == 0 {
                if rand.Intn(2) == 1 {
                    avatar.SetRGBA(x, y, pixelColor)
                }
            }
        }
    }
    
    // Mirror left half to right half
    for y := 0; y < thumbSize; y++ {
        for x := 0; x < thumbSize; x++ {
            if x < thumbSize/2 {
                avatar.Set(thumbSize-x-1, y, avatar.At(x, y))
            }
        }
    }
    
    // Mirror top to bottom
    for y := 0; y < thumbSize; y++ {
        for x := 0; x < thumbSize; x++ {
            if y < thumbSize/2 {
                avatar.Set(x, thumbSize-y-1, avatar.At(x, y))
            }
        }
    }
}

func paintBG(avatar *image.RGBA, bgColor color.RGBA) {
    for y := 0; y < thumbSize; y++ {
        for x := 0; x < thumbSize; x++ {
            avatar.SetRGBA(x, y, bgColor)
        }
    }
}

func calcPixelColor(nameBytes []byte) (pixelColor color.RGBA) {
    pixelColor.A = 255
    
    var mutator = byte((len(nameBytes) * 4))
    
    pixelColor.R = nameBytes[0] * mutator
    pixelColor.G = nameBytes[1] * mutator
    pixelColor.B = nameBytes[2] * mutator
    
    return
}

func calcBGColor(nameBytes []byte) (bgColor color.RGBA) {
    bgColor.A = 255
    
    var mutator = byte((len(nameBytes) * 2))
    
    bgColor.R = nameBytes[0] * mutator
    bgColor.G = nameBytes[1] * mutator
    bgColor.B = nameBytes[2] * mutator
    
    return
}