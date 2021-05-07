package main

import (
	"fmt"
	"litrocket/common"
	"litrocket/controler"
	"litrocket/model"
	"litrocket/router"
	"litrocket/utils/handlelog"
	"net"

	"gopkg.in/ini.v1"
)

// Config File
const Configfile = "config/rocket.ini"

func main() {
	// Image.
	StrImg()

	// Read Config.
	ReadConfig()

	// Init Servers.
	InitServers()

	// Init DataBase.
	model.InitDb()

	// Init Router.
	router.InitRouter()

	handlelog.Handlelog("INFO", "Litrocket is running......")

	// Start Running...
	controler.HandleConn()
}

// Read Config From Configfile.
func ReadConfig() {
	f, err := ini.Load(Configfile)
	if err != nil {
		handlelog.Handlelog("FATAl", "配置文件加载错误!")
	}

	server := f.Section("server")
	db := f.Section("database")
	encry := f.Section("security")

	common.AppMode = server.Key("AppMode").MustString("debug")

	common.RequestAddr = server.Key("RequestServer").MustString("")
	common.ResponseAddr = server.Key("ResponseServer").MustString("")
	common.ChatAddr = server.Key("ChatServer").MustString("")
	common.FileAddr = server.Key("FileServer").MustString("")
	common.VideoAddr = server.Key("VideoServer").MustString("")

	common.DbHost = db.Key("DbHost").MustString("localhost")
	common.DbUser = db.Key("DbUser").MustString("root")
	common.DbPort = db.Key("DbPort").MustString("3306")
	common.DbName = db.Key("DbName").MustString("litrocket")
	common.DbPass = db.Key("DbPass").MustString("")

	common.AESKEY = encry.Key("aeskey").MustString("321423u9y8d2fwfl")
}

// Init Servers (Create Listeners By Config) .
func InitServers() {
	var err error

	if common.RequestListener, err = net.Listen("tcp", common.RequestAddr); err != nil {
		handlelog.Handlelog("FATAL", "Init RequestServer Error"+err.Error())
	}
	handlelog.Handlelog("INFO", "Init RequestServer ok "+common.RequestAddr)

	if common.ResponseListener, err = net.Listen("tcp", common.ResponseAddr); err != nil {
		handlelog.Handlelog("FATAL", "Init ResponseServer Error"+err.Error())
	}
	handlelog.Handlelog("INFO", "Init ResponseServer ok"+common.ResponseAddr)

	if common.ChatListener, err = net.Listen("tcp", common.ChatAddr); err != nil {
		handlelog.Handlelog("FATAL", "Init ChatServer Error"+err.Error())
	}
	handlelog.Handlelog("INFO", "Init ChatServer ok "+common.ChatAddr)

	if common.FileListener, err = net.Listen("tcp", common.FileAddr); err != nil {
		handlelog.Handlelog("FATAL", "Init FileServer Error"+err.Error())
	}
	handlelog.Handlelog("INFO", "Init FileServer ok "+common.FileAddr)

	if common.VideoListener, err = net.Listen("tcp", common.VideoAddr); err != nil {
		handlelog.Handlelog("FATAL", "Init VideoServer Error"+err.Error())
	}
	handlelog.Handlelog("INFO", "Init VideoServer ok "+common.VideoAddr)
}

// Print String Image.
func StrImg() {
	str := `
           .,:,,,                                        .::,,,::.
         .::::,,;;,                                  .,;;:,,....:i:
         :i,.::::,;i:.      ....,,:::::::::,....   .;i:,.  ......;i.
         :;..:::;::::i;,,:::;:,,,,,,,,,,..,.,,:::iri:. .,:irsr:,.;i.
         ;;..,::::;;;;ri,,,.                    ..,,:;s1s1ssrr;,.;r,
         :;. ,::;ii;:,     . ...................     .;iirri;;;,,;i,
         ,i. .;ri:.   ... ............................  .,,:;:,,,;i:
         :s,.;r:... ....................................... .::;::s;
         ,1r::. .............,,,.,,:,,........................,;iir;
         ,s;...........     ..::.,;:,,.          ...............,;1s
        :i,..,.              .,:,,::,.          .......... .......;1,
       ir,....:rrssr;:,       ,,.,::.     .r5S9989398G95hr;. ....,.:s,
      ;r,..,s9855513XHAG3i   .,,,,,,,.  ,S931,.,,.;s;s&BHHA8s.,..,..:r:
     :r;..rGGh,  :SAG;;G@BS:.,,,,,,,,,.r83:      hHH1sXMBHHHM3..,,,,.ir.
    ,si,.1GS,   sBMAAX&MBMB5,,,,,,:,,.:&8       3@HXHBMBHBBH#X,.,,,,,,rr
    ;1:,,SH:   .A@&&B#&8H#BS,,,,,,,,,.,5XS,     3@MHABM&59M#As..,,,,:,is,
   .rr,,,;9&1   hBHHBB&8AMGr,,,,,,,,,,,:h&&9s;   r9&BMHBHMB9:  . .,,,,;ri.
   :1:....:5&XSi;r8BMBHHA9r:,......,,,,:ii19GG88899XHHH&GSr.      ...,:rs.
   ;s.     .:sS8G8GG889hi.        ....,,:;:,.:irssrriii:,.        ...,,i1,
   ;1,         ..,....,,isssi;,        .,,.                      ....,.i1,
   ;h:               i9HHBMBBHAX9:         .                     ...,,,rs,
   ,1i..            :A#MBBBBMHB##s                             ....,,,;si.
   .r1,..        ,..;3BMBBBHBB#Bh.     ..                    ....,,,,,i1;
    :h;..       .,..;,1XBMMMMBXs,.,, .. :: ,.               ....,,,,,,ss.
     ih: ..    .;;;, ;;:s58A3i,..    ,. ,.:,,.             ...,,,,,:,s1,
     .s1,....   .,;sh,  ,iSAXs;.    ,.  ,,.i85            ...,,,,,,:i1;
      .rh: ...     rXG9XBBM#M#MHAX3hss13&&HHXr         .....,,,,,,,ih;
       .s5: .....    i598X&&A&AAAAAA&XG851r:       ........,,,,:,,sh;
       . ihr, ...  .         ..                    ........,,,,,;11:.
          ,s1i. ...  ..,,,..,,,.,,.,,.,..       ........,,.,,.;s5i.
           .:s1r,......................       ..............;shs,
           . .:shr:.  ....                 ..............,ishs.
               .,issr;,... ...........................,is1s;.
                  .,is1si;:,....................,:;ir1sr;,
                     ..:isssssrrii;::::::;;iirsssssr;:..
                          .,::iiirsssssssssrri;;:.


                                  ... C...
                                C          C
                              C              C
                            C                  C
                           C                     C
                            LTTTTTTTTTTTTTTTTTTTT       
                           /L R R R    EEEEEEEE T\
                           /L R    R   E        T\
                           /L R     R  E        T\
                           /L R    R   E        T\
                           /L R R R    EEEEEEEE T\
                           /L R  R     E        T\
                           /L R   R    E        T\
                           /L R    R   E        T\
                           /L R     R  EEEEEEEE T\
                           /L                   T\
                           /L         O         T\
                           /L      O     O      T\
                           /L    O         O    T\
                           /L  O             O  T\
                           /L    O         O    T\
                           /L      O     O      T\
                           /L         O         T\
                           /L                   T\
                           /LLLLLLLLLLLLLLLLLLLLT\
                           .KKKKKKKKKKKKKKKKKKKKK.
                                     .K.
                                   .K. .K.
                                  .K.   .K.
                                 .K.     .K.
                                .K.       .K.
                               .K.         .K.
                              .K.           .K.
                             .K.             .K.

 __        ______  ________  _______    ______    ______   __    __  ________  ________ 
/  |      /      |/        |/       \  /      \  /      \ /  |  /  |/        |/        |
$$ |      $$$$$$/ $$$$$$$$/ $$$$$$$  |/$$$$$$  |/$$$$$$  |$$ | /$$/ $$$$$$$$/ $$$$$$$$/ 
$$ |        $$ |     $$ |   $$ |__$$ |$$ |  $$ |$$ |  $$/ $$ |/$$/  $$ |__       $$ |   
$$ |        $$ |     $$ |   $$    $$< $$ |  $$ |$$ |      $$  $$<   $$    |      $$ |   
$$ |        $$ |     $$ |   $$$$$$$  |$$ |  $$ |$$ |   __ $$$$$  \  $$$$$/       $$ |   
$$ |_____  _$$ |_    $$ |   $$ |  $$ |$$ \__$$ |$$ \__/  |$$ |$$  \ $$ |_____    $$ |   
$$       |/ $$   |   $$ |   $$ |  $$ |$$    $$/ $$    $$/ $$ | $$  |$$       |   $$ |   
$$$$$$$$/ $$$$$$/    $$/    $$/   $$/  $$$$$$/   $$$$$$/  $$/   $$/ $$$$$$$$/    $$/    
`
	fmt.Println(str)
}
