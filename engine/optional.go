 package engine

 type Optional[T any] interface {
   HasValue() bool
   Value() T
 }

 // Here methods if it has any value 
 type SomeValue[T any] struct {
   value T 
 }

 func (s *SomeValue[T]) HasValue() bool {
   return true 
 }

 func (s *SomeValue[T]) Value() T {
   return s.value 
 }

 //Here methods if it has no value
 type NoValue[T any] struct {}

 func(n NoValue[T]) HasValue() bool {
   return false
 }

 func (n NoValue[T]) Value() T {
   return *new(T)
 }

// Example usage:
//  - Create an optional string with a value:
//     optString := optional.SomeValue[string]{value: "Hello"}
//  - Create an optional string without a value:
//     optString = optional.NoValue[string]{}
//  - Check for the presence of a value:
//     if optString.HasValue() {
//         // Value is present, use it
//         fmt.Println("The string value is:", optString.Value())
//     } else {
//         // Value is absent, handle accordingly
//         fmt.Println("No string value present")
//     }
// how to use it in a struct
//  email Optional[string] // Optional string field
