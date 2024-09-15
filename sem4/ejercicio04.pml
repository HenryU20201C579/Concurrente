//inicio de procesos concurrente
byte n

proctype P(byte id; byte inc){
    byte temp
    temp=n + inc
    n=temp
    printf("Proceso P(%d), n=%d\n",id,n)
}

init{
    n=0
    atomic{
        run P(1,10)
        run P(2,15)
    }
    (_nr_pr == 1) -> printf("El valor final de n es %d\n",n)
}