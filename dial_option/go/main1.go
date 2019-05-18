package main

import "fmt"

type Option interface {
    apply(*A)
}

type funcOption struct {
    f func(*A)
}

func (fopt *funcOption) apply(a *A) {
    fopt.f(a)
}

func WithArg1(arg1 string) Option {
    return &funcOption{func(a *A) {
        a.arg1 = arg1
    }}
}

func WithArg2(arg2 int) Option {
    return &funcOption{func(a *A) {
        a.arg2 = arg2
    }}
}

type A struct {
    arg1 string
    arg2 int
}

func NewA(opts ...Option) *A {
    a := &A{}
    for _, opt := range opts {
        opt.apply(a)
    }
    return a
}

func (a *A) String() string {
    return fmt.Sprintf("&A{arg1: %s, arg2: %d}", a.arg1, a.arg2)
}

func main() {
    a := NewA(WithArg1("1"), WithArg2(1))
    fmt.Println(a)
}