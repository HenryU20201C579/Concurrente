//programa concurrente básico
int n=0

proctype P(){
    int k1=1
    n=k1
}

proctype Q(){
    int k2=2
    n=k2
}

init{
    atomic{
        run P()
        run Q()
    }
    (_nr_pr == 1) -> printf("El valor final de n es %d\n",n)
    //condición = código correcto
    assert(n==1)  //válido
}