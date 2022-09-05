# add-key

So recently-ish debian deprecated using `apt-key add -` in favor of adding individual signed keys. This is all well and good but [this|https://unix.stackexchange.com/questions/332672/how-to-add-a-third-party-repo-and-key-in-debian] process is a pain in the butt. So I created this simple little utility to make it easier.

The following is an example to install the spotify repo https://www.spotify.com/us/download/linux/
```
kellen@chewbacca:~/pkg/add-key$ sudo ./add-key spotify --gpg=https://download.spotify.com/debian/pubkey_5E3C45D7B312C643.gpg --suite=stable --components=non-free --uri=http://repository.spotify.com
Wrote /usr/share/keyrings/spotify.gpg
Wrote /etc/apt/sources.list.d/spotify.sources
```
