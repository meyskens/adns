A(dobe) DNS Server
==================
A(dobe) DNS Server or ADNS for short is a DNS-over-HTTPS proxy that only allows DNS to Adobe domains to go through. It is designed to be used for breaking internet access except for Adobe CC licensing. A requirement needed for when using Adobe software during exams.

## How to use
1) Block port 53 on the network (we use DoH to bypass this)
2) Change your local DNS server to `127.0.0.1`
3) Run adns, in Windows just open `adns.exe`
4) Open Adobe CC
*Note: currently the app logs blocked domain names, this is for debugging purposes please ignore this*