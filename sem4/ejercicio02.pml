//procesos dependientes
proctype hijo(){
    printf("Hijo (%d)\n",_pid)
}

active proctype padre(){
    do
    :: (_nr_pr == 1) -> 
                        run hijo()
    od
}