package sci


// Interface conformity confirmations
var _ Unit = &BaseUnit{}
var _ Unit = &DerivedUnit{}
var _ Unit = &NilUnit{}
var _ Unit = &DivUnit{}
var _ Unit = &MulUnit{}

var _ error = &MagnitudeError{}
var _ error = &ExpToBigError{}
var _ error = &BaseUnitAlreadyDefinedError{}

