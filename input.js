const list = [1, 2, 3, 4]
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


console.log(twoSum(list, 6))
