function twoSum(arr, target) {
        let left = 0
        let right = arr.length

        while (left < right) {
                let sum = arr[left] + arr[right]
                if (sum == target) {
                        return sum 
                } else {
                        left++ 
                        right--
                }

        }

        return 0
}

const arr = [1, 2, 3, 4]
console.log(twoSum(arr, 6))
