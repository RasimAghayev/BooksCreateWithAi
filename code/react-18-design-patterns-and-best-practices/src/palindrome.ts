// Chapter 2 — Converting JavaScript to TypeScript: palindromic check (page 46-47)

// Original JS:
// function isPalindrome(word) {
//   const lowerCaseWord = word.toLowerCase()
//   const reversedWord = lowerCaseWord.split('').reverse().join('')
//   return lowerCaseWord === reversedWord
// }

// TS: add parameter + return types (any valid JS is valid TS)
function isPalindrome(word: string): boolean {
  const lowerCaseWord = word.toLowerCase()
  const reversedWord = lowerCaseWord.split('').reverse().join('')
  return lowerCaseWord === reversedWord
}

// Usage / compile-time type checking (Figure 2.1):
console.log(isPalindrome('Level'))  // true
console.log(isPalindrome('Anna'))   // true
console.log(isPalindrome('Carlos')) // false
// Passing a non-string fails at compile time:
// console.log(isPalindrome(101))    // TS Error: number not assignable to string
// console.log(isPalindrome(true))   // TS Error
// console.log(isPalindrome(false))  // TS Error
