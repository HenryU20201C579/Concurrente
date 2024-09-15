//contador
#define N 10

active proctype P(){
    int sum=0
    byte i=1
    do
    :: i>N -> break
    :: else ->
            sum=sum+i
            i++
    od
    printf("La suma del número %d es %d\n",N,sum)
}