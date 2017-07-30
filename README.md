# umirus
umirus generates programs which iteratively copy included datas (like computer virus. Of course, its not virus.)

## Motivation
I would like to see a beautiful sea and feel at ease... so I decided to fill up my desktop with "うみ" pngs automatically. (umi[うみ/海] means sea/ocean in Japanese...)

## Installation

    go get github.com/paka3m/umirus
    umirus

![umirusimplementation](https://user-images.githubusercontent.com/11421401/28755586-2182deda-7599-11e7-98d1-baba586de3fe.png)

## how to reprogram as you like... 

    git clone git@github.com:paka3m/umirus.git
    cd umirus
    
You have to copy the (image)-files into ```umirus/assets``` folder. 

In this case, ```umirus/assets/umi``` is a place of example cute images related to sea, which you can download freely from http://www.irasutoya.com/2013/11/blog-post_3484.html.

    ls ./assets/umi
    >   animal_jinbeizame.png
        ocean_kamome.png
        ocean_kurage.png
        ocean_ukiwa.png
        animal_shachi_killer_whale.png ...

Next, you have to include these files to the binary by using ```go-bindata```.

    go get -u github.com/jteeuwen/go-bindata/...
    go-bindata -ignore .DS_Store assets/umi

and build it with ldflags.

    go build -ldflags "-X main.filename=うみうみFile -X main.dir=assets/umi -X main.wait=1s" -o umirusd

now ```umirusd``` program is generated.

while executing this program, images will be generated.

![GeneratedManyFiles](https://user-images.githubusercontent.com/11421401/28755340-54298966-7593-11e7-904e-e9d27e5e02f4.png)

![GeneratedMuchManyFiles](https://user-images.githubusercontent.com/11421401/28755328-01d16f94-7593-11e7-9048-31378db6f543.png)

## How to stop (Unix)
Don't be afraid. It's very easy! All you have to do is to command

    pgrep -f 'umirusd' | xargs kill

or use GUI such as activity monitor(macOS).

# Caution
this application is designed for experimental use.

# Link
いらすとや様 http://www.irasutoya.com/

# Author
paka3m

# Licence
MIT