I have used sigma23[140.158.128.21] for server and printer
and [140.158.128.15],    [140.158.128.12],    [140.158.130.213],    [140.158.131.163] for fork and philosopher
all the corresponding files are placed in all servers within a directory 'project2'
step 1: Run the server.go
	command : go run server.go -n 5 -host "140.158.128.21:1111"    ----- [ I have used 1111 as port number here]

step 2: To display the execution of the process run printer.go		      [ I have used the same host 'sigma23' for the printer as well]
	command : go run printer.go -n 5 -host "140.158.128.21:2222"    ----- [ I have used 2222 as port number here]

step 3:
Once the server and display server are started  run the fork using the following commands
    cd project2
	go run fork.go -id <Unique Id>  -host<Ip address> -manager <Ip address of the manager server>
    ex:
	go run fork.go -id 0 -host "localhost:11111" -manager "140.158.128.21:1111"
	go run fork.go -id 1 -host "localhost:11111" -manager "140.158.128.21:1111"
   	go run fork.go -id 2 -host "localhost:11111" -manager "140.158.128.21:1111"
  	go run fork.go -id 3 -host "localhost:11111" -manager "140.158.128.21:1111"
   	go run fork.go -id 4 -host "localhost:11111" -manager "140.158.128.21:1111"
    
step 4:    
Now start 5 philosopher service by executing following commands
   cd project2
	go run philosopher.go -id <Unique Id> -n <No.of Philo's/Forks > -host <Ip address> -manager <Ip address of manager server> -printer <Ip address of display server>
    Ex:-
   	go run philosopher.go -id 0 -n 5 -host "localhost:11111" -manager "140.158.128.21:1111" -printer "140.158.128.21:2222"
   	go run philosopher.go -id 1 -n 5 -host "localhost:11111" -manager "140.158.128.21:1111" -printer "140.158.128.21:2222"
   	go run philosopher.go -id 2 -n 5 -host "localhost:11111" -manager "140.158.128.21:1111" -printer "140.158.128.21:2222"
   	go run philosopher.go -id 3 -n 5 -host "localhost:11111" -manager "140.158.128.21:1111" -printer "140.158.128.21:2222" 
   	go run philosopher.go -id 4 -n 5 -host "localhost:11111" -manager "140.158.128.21:1111" -printer "140.158.128.21:2222" 
