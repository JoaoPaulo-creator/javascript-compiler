# Steps to follow in my mind

## Lexer

### First part
- create some structure with field members
- create a function that reads the characters from input
- build this lexer

the next steps is the tokenization per se
I need to create the function that runs the
lexical analysis, but a lot of helper functions should be written/created to perform the tokenization

### Second part - helper functions

- the first function is to skip whitespaces
- then a function to peek the current character
- create function to check if is digit or number
- create a function to build the Token struct


... for now let's build the function NewToken, to start the tokenization


```js


const list = [1, 2, 3, 4]
function twoSum(arr, target) {
        let left = 0
        let right = arr.length

        while (left < right) {
                sum = arr[left] + arr[right]
                if (sum == target) {
                        return sum 
                } else {
                        left++ 
                        right--
                }

        }

        return 0
}
```


