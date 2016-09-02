A go library for scientific unit calculations, powering a scientific spreadsheet application.  This library is inspired by ruby-unit and allows for code like:

```go

  // creating values
  x := &sci.Value{M: "10.001", U: si.Millimeter}
  y := &sci.Value{M: "10.001", U: si.Millimeter}
  z := &sci.Value{}

  // math
  err := z.Add(x, y)

  // units are composite values, algebraic values

  // acceleration (m/s^2)
  Acceleration = sci.DivUnit{
    N: si.Meter,
    D: sci.MulUnit{ units.Second, units.Second },
  }
   

  // Defining constants
  var (
    EarthGravity = si.MustParseValue("9.807 m/s^2")
  )
```

This spreadsheet application could allow for "input cells" and "output cells" which exposes the sheet as an API to others allowing for collaboration.