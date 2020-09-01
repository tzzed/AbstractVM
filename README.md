A third year Epitech project written in C++. I decided to rewrite it in Go


The goal was to develop a virtual machine capable of basic arithmetic operations (+,-,*,/,%).

The way to communicate with the machine is an assembly language whose the grammar is specified in grammar.txt

There is an example of code in example.avm

### AbstractVM

### Install

Install AbstractVM interpreter

```
git clone github.com/tzzed/AbstractVM.git
cd AbstractVM
make
```

### Usage 

### REPL
```
$>avm
Abstract VM
Enter ".help" for usage hints.
avm>
```

### File interpreter

```
$>avm f.avm
{42 int32}
{42.42 double}
{119.55 float}
$>








