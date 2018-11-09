# Goals
Your personal goals register and tracker.

Simple but powerful, written in __Go__.

![Creating a new goal ](https://im2.ezgif.com/tmp/ezgif-2-30b26bc81bdf.gif)

## Installation
```
cd $GOPATH
go get github.com/giogiglio/goals
go install goals
```

Make sure to have `$GOPATH/bin` in your `$PATH` so that you can execute `goals` from everywhere in your terminal.

### Dependencies
These packages are required in order to run `goals`.
- [go-sqlite3](https://github.com/mattn/go-sqlite3) sqlite3 driver for Go.
- [Survey](https://github.com/AlecAivazis/survey) for interactive prompts.

## Usage
```
goals [ [-new | -edit | -remove] goal | progress ]
```

### Constraints:
- Goals names must be unique and must not exceed 20 characters.
- Goals date format is `dd/mm/yyyy`. You can also use `today` and `yesterday`.
- Goals notes can be empty and can't exceed 50 characters.
- Progresses value must be a number between 0 and 100 included.
- As for goals notes, the same applies to progresses ones.
