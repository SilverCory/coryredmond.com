package site

import (
	"hash/fnv"
	"math/rand"
	"strconv"
	"time"
)

var ASCII = []string{
	`
		     ******       ******
		   **********   **********
		 ************* *************
		*****************************
		*****************************
		*****************************
		 ***************************
		   ***********************
		     *******************
		       ***************
		         ***********
		           *******
		             ***
		              *             You can check out any time you want, but you can never leave.`,

	`
		        @@@@@@           @@@@@@
		      @@@@@@@@@@       @@@@@@@@@@
		    @@@@@@@@@@@@@@   @@@@@@@@@@@@@@
		  @@@@@@@@@@@@@@@@@ @@@@@@@@@@@@@@@@@
		 @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
		@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
		@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
		@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
		 @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
		  @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
		   @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
		    @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
		      @@@@@@@@@@@@@@@@@@@@@@@@@@@
		        @@@@@@@@@@@@@@@@@@@@@@@
		          @@@@@@@@@@@@@@@@@@@
		            @@@@@@@@@@@@@@@
		              @@@@@@@@@@@
		                @@@@@@@
		                  @@@
		                   @       Can't go home alone again, need someone to numb the pain`,
	`
		  ff     ff
		 f   f   f   f
		f     f f     f
		f      f      f
		f             f
		 f           f
		  f         f
		   f       f
		    f     f
		     f   f
		      f f
		       f        Press f to pay respects`,
	`
		         ,ae,
		       ,88888e
		,a888b.9888888i
		888888888888888
		88888888888888Y
		'8888888888888'
		  "S888888888"
		    "7888888Y
		       "e88j
		         "Y     What is love? Baby don't hurt me.`,
	`
        ▒▒▒░░░░░░░░░░▄▐░░░░
        ▒░░░░░░▄▄▄░░▄██▄░░░
        ░░░░░░▐▀█▀▌░░░░▀█▄░
        ░░░░░░▐█▄█▌░░░░░░▀█▄
        ░░░░░░░▀▄▀░░░▄▄▄▄▄▀▀
        ░░░░░▄▄▄██▀▀▀▀░░░░░
        ░░░░█▀▄▄▄█░▀▀░░░░░░
        ░░░░▌░▄▄▄▐▌▀▀▀░░░░░
        ░▄░▐░░░▄▄░█░▀▀░░░░░
        ░▀█▌░░░▄░▀█▀░▀░░░░░
        ░░░░░░░░▄▄▐▌▄▄░░░░░
        ░░░░░░░░▀███▀█░▄░░░
        ░░░░░░░▐▌▀▄▀▄▀▐▄░░░
        ░░░░░░░▐▀░░░░░░▐▌░░
        ░░░░░░░█░░░░░░░░█░░
        ░░░░░░▐▌░░░░░░░░░█░      Spooky scary skeletons!`,
	`
		                         )
		                         (
		           (            ( )
		           )             "
		          ( )
		           "    |          |
		       |       ((          ))
		 |     ))      ))   )     //
		 ))   ( )     / /   (    ( (     |
		 .(    \ \   ( (   ( )   )  )    ))
		( (     ) '.'   '.  "  .'   /   ( \
		 ) \  .'          '._.'     '._  ) )
		 :  '' _.oooooo._   _.oooooo._ ''  /
		  )  .odOOOOOOOObo.odOOOOOOOObo.   |
		 /  dOOOY dOOOOOOOOOOOOOOOOOOOOOb  (
		   OOOY dOOOOOOOOOOOOOOOOOOOOOOOOO  \
		  dOOY dOOOOOOOOOOOOOOOOOOOOOOOOOOb
		  OOO dOOOOOOOOOOOOOOOOOOOOOOOOOOOO
		  OOO YOOOOOOOOOOOOOOOOOOOOOOOOOOOO
		  YOOb OOOOOOOOOOOOOOOOOOOOOOOOOOOY
		   YOObdOOOOOOOOOOOOOOOOOOOOOOOOOY
		    YOOOOOOOOOOOOOOOOOOOOOOOOOxXY
		     "YOOOOOOOOOOOOOOOOOOOOOxXY"
		       "YOOOOOOOOOOOOOOOOOAoS"
		         YOOOOOOOOOOOOOxXXXY
		          "YOOOOOOOOOxXXXY"
		             "YOOOXXXXY"
		                ""Y""       We're building it up, to burn it back down.'`,
	`
		888                            888    
		888                            888    
		888                            888    
		88888b.  .d88b.  8888b. 888d888888888 
		888 "88bd8P  Y8b    "88b888P"  888    
		888  88888888888.d888888888    888    
		888  888Y8b.    888  888888    Y88b.  
		888  888 "Y8888 "Y888888888     "Y888    I couldn't find any more good ascii art hearts..'`,
	`
		 _________________
		|# :           : #|
		|  :           :  |
		|  :           :  |
		|  :           :  |
		|  :___________:  |
		|     _________   |
		|    | __      |  |
		|    ||  |     |  |
		\____||__|_____|__|    I was only a kid when I used these.`,
	`
		BBBBBBBBBBBBBBBBBBBBBBBBBBB
		BMB---------------------B B
		BBB---------------------BBB
		BBB---------------------BBB
		BBB---------------------BBB
		BBB---------------------BBB
		BBB---------------------BBB
		BBBBBBBBBBBBBBBBBBBBBBBBBBB
		BBBBB++++++++++++++++BBBBBB
		BBBBB++BBBBB+++++++++BBBBBB
		BBBBB++BBBBB+++++++++BBBBBB
		BBBBB++BBBBB+++++++++BBBBBB
		BBBBB++++++++++++++++BBBBBB   Cocking the floppy disk like a gun, chk chk!`,
	`
		 ___________________________.
		|;;|                     |;;||
		|[]|---------------------|[]||
		|;;|                     |;;||
		|;;|                     |;;||
		|;;|                     |;;||
		|;;|                     |;;||
		|;;|                     |;;||
		|;;|                     |;;||
		|;;|_____________________|;;||
		|;;;;;;;;;;;;;;;;;;;;;;;;;;;||
		|;;;;;;_______________ ;;;;;||
		|;;;;;|  ___          |;;;;;||
		|;;;;;| |;;;|         |;;;;;||
		|;;;;;| |;;;|         |;;;;;||
		|;;;;;| |;;;|         |;;;;;||
		|;;;;;| |;;;|         |;;;;;||
		|;;;;;| |___|         |;;;;;||
		\_____|_______________|_____||
		 ~~~~~^^^^^^^^^^^^^^^^^~~~~~~	Floppy disks are by far my favourite storage medium. Please bring them back..`,
	`
		       !
		       !
		       ^
		      / \
		     /___\
		    |=   =|
		    |     |
		    |     |
		    |     |
		    |     |
		    |     |
		    |     |
		    |     |
		    |     |
		    |     |
		   /|##!##|\
		  / |##!##| \
		 /  |##!##|  \
		|  / ^ | ^ \  |
		| /  ( | )  \ |
		|/   ( | )   \|
		    ((   ))
		   ((  :  ))
		   ((  :  ))
		    ((   ))
		     (( ))
		      ( )
		       .
		       .
		       .    Space, it's cool.'`,
	`
             ROFL:ROFL:ROFL:ROFL
                  ___ ^_____
         L    ___/         [ ]
        LOL===_
         L     \_____________]
                 ___I______I__      rofl`,
	`
               ▄▀▀▀▀▀▀▀▀▀▀▄▄
            ▄▀▀             ▀▄
          ▄▀                  ▀▄
          █                     ▀▄
         ▐▌        ▄▄▄▄▄▄▄       ▐▌
         █           ▄▄▄▄  ▀▀▀▀▀  █
        ▐▌       ▀▀▀▀     ▀▀▀▀▀   ▐▌
        █         ▄▄▀▀▀▀▀    ▀▀▀▀▄ █
        █                ▀   ▐     ▐▌
        ▐▌         ▐▀▀██▄      ▄▄▄ ▐▌
         █           ▀▀▀      ▀▀██  █
         ▐▌    ▄             ▌      █
          ▐▌  ▐              ▀▄     █
           █   ▌        ▐▀    ▄▀   ▐▌
           ▐▌  ▀▄        ▀ ▀ ▀▀   ▄▀
           ▐▌  ▐▀▄                █
           ▐▌   ▌ ▀▄    ▀▀▀▀▀▀   █
           █   ▀    ▀▄          ▄▀
          ▐▌          ▀▄      ▄▀
         ▄▀   ▄▀        ▀▀▀▀█▀
        ▀   ▄▀          ▀   ▀▀▀▀▄▄▄▄▄       Do I need to say anything for this one...`,
	`
        █                                                                             █
        █░▀▓█▀ ███▄    █   ██████ ▓█████ ▄████▄  █    ██  ██▀███   ██░████████▓▓██   ██
        ▓ ░██  ██ ▀█░  █ ▒██    ▒ ▓█   ▀▒██▀ ▀█  ██  ▓██▒▓██ ▒ ██▒▓██░▓  ██▒ ▓▒ ▒██  ██
        ▓ ░██ ▓██ ░▀█ ██▒░ ▓██▄   ▒███  ▒▓█    ▄▓██  ▒██░▓██ ░▄█ ▒▒██ ▒ ▓██░ ▒░  ▒██ ██
        ▒ ░█▓░▓██▒ ░▐▌██▒  ▒   ██▒▒██  ▄▒▓▓▄ ▄██▓▓█  ░██░▒██▀▀█▄  ░██░░ ▓██▓ ░   ░ ▐██▓
        ░░▄██▄▒██░   ▓██░▒██████▒▒░▒████▒ ▓███▀ ▒▒█████▓ ░██  ▒██▒░██░  ▒██▒ ░   ░ ██▒▓
        ░░ ░░░░ ▒    ▒     ▒▓▒ ▒     ▒░   ░▒ ▒  ░  ▒ ▒ ▒ ░ ▒  ░▒ ░░▓    ▒ ░░      ██▒░▓
        ░  ░                ░                                  ░                ▓██   ▒      Does it count if it's voluntary?'`,
}

func GetAscii(ip string) string {
	rand.Seed(FNV32a(ip + strconv.Itoa(time.Now().Day())))
	return ASCII[rand.Int()%len(ASCII)]
}

func FNV32a(text string) int64 {
	algorithm := fnv.New64a()
	algorithm.Write([]byte(text))
	return int64(algorithm.Sum64())
}
