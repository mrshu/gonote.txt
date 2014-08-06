CC=go
INSTALL=install -c
DEST=/usr/bin

gonote: gonote.go
	$(CC) build ./gonote.go

install: gonote 
	$(INSTALL) gonote $(DEST)/gonote

clean:
	rm -rf gonote 
