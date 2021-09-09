This tool may be of some use for admins responsible for managing one or more Kea servers (DHCP).

To build:\
go mod init\
go get <external modules>\
go build -o keatool .

You will need to import:

        "fmt"
        "bufio"
        "os"
        "strings"
        "bytes"
        "encoding/binary"
        "net"
        "database/sql"
        "sort"
        "log"
        "github.com/go-sql-driver/mysql"
        "github.com/DavidGamba/go-getoptions"
        
some of which are external dependencies.  

Use Case:

This is a goofy little CLI tool I wrote in Go to do some basic administration tasks to the mysql backend I use for kea (1.6.2).

The file, ktool.cfg will hold the hostname, username, and password for accessing your mysql instance.

This builds OK on Windows, MacOS X 10.14, and CentOS 7.  On Windows, it will require you to run the tool under a 'real' Unix/Linux shell.  
MobaXterm works in my testing, but there are other options!

Notes:

This was written to use TLS between the client and the mysql server.  You will need to configure a TLS cert + key on the mysql server and make sure the certificate chain can be traversed by the client system.  We use this with real certs and it has not been tested with self-signed certificates.

        
