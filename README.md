# Cryptofinder

## Find all directories and files that cryptowall has taken hostage

Im not sure if this will work on every PC

but on my PC I noticed the first 16 bytes of every encrypted file was the same...

pre built binaries included in the dist directory


**Usage**
```
cryptofinder <start directory> [clean]
```

**Ex**
scan all files in the c:\ recursively
```
cryptofinder c:\
```

scan all files in the c:\ recursively and delete all the HELP_DECRYPT files while searching
```
cryptofinder c:\ clean
```

scan all files in the /home/pborges directory recursively
```
cryptofinder /home/pborges
```

scan all files in the /home/pborges directory recursively and delete all the HELP_DECRYPT files while searching
```
cryptofinder /home/pborges clean
```
