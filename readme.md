Razor - A BitTorrent covered shell
* this project should be used for educational use only

Main features:
* BitTorrent-only connectivity
* Remote command execution
* AES256 encryption
* File upload
* Basic tunneling capabilities 

How to install:

    1. Create an info hash (any sha1 will do)
    2. Decide on a remote port to listen on
    3. Create a 256 bit long encryption key
    4. Compile as exe
    5. Run PyRazor using your key and hash
    6. Connect and have fun.
    
Minor disclosures:
* The encryption is done poorly, it's only there so that
the content of commands and uploaded files will not be in clear text

* To a BitTorrent expert the communication will look odd
but anyone who didn't read the protocol documentation would not find a difference
     
    
