package task



#GenInt32: {
	$dolusTask: "GenInt32"
    min: int32 | *0
    max: int32 &>=min | *10
}

#Sequentional: {
    $dolusTask: "Sequentional"
    start: int | *0
}

