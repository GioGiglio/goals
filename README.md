# Goals
Your personal goals register and tracker.

Simple but powerful, written in __Go__.

![Creating a new goal ](https://image.ibb.co/gHFKVA/ezgif-com-gif-maker-2.gif)

## Installation
```
go get github.com/giogiglio/goals
cd $GOPATH/src
go install github.com/giogiglio/goals
```

Make sure to have `$GOPATH/bin` in your `$PATH` so that you can execute `goals` from everywhere in your terminal.

### Dependencies
These modules are required.
- [go-sqlite3](https://github.com/mattn/go-sqlite3) sqlite3 driver for Go.
- [Survey](https://github.com/AlecAivazis/survey) for interactive prompts.

## Usage
```
goals [ [-new | -edit | -remove] goal | progress ] [ -help ]

goals -new [ goal | progress ]
  Lets you create a new goal, or a new progress for an exisisting goal.
  
goals -edit [ goal | progress ]
  Lets you modify an existing goal or a progress.
  
goals -remove [ goal | progress ]
  Lets you remove an existing goal or progress.
  
goals -help
  Prints an help message.
```

### Constraints:
- Goal name must be unique and must not exceed 20 characters.
- Goal and progress date format is `dd/mm/yyyy`. You can also use `today` and `yesterday`.
- Goal and progress note can be empty and can't exceed 50 characters.
- Progress value must be a number between 0 and 100, both included.
