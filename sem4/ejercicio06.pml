mtype = {red,yellow,green}
mtype foco = green

active proctype P(){
    do
    :: if
        :: foco == red -> foco=green
        :: foco == yellow -> foco=red
        :: foco == green -> foco=yellow
       fi
       printf("El foco es ahora %e\n", foco)
    od
}