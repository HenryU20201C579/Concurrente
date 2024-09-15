//comparativo No deterministico
active proctype P(){
    int a=5, b=5
    int max
    int branch
    if
    :: a>=b -> max=a
              branch=1
    :: b>=a -> max=b
              branch=2
    fi
    printf("El máximo entre %d and %d es %d en el branch %d\n",a,b,max,branch)
}