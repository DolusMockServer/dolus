package task



#GenInt32: {
	$dolus: task: "GenInt32"
    min: int32 | *0
    max: int32 &>=min | *10
}

#Sequentional: {
    $dolus: task: "Sequentional"
    start: int | *0
}

