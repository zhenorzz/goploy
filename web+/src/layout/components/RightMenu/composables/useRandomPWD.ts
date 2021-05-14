import { reactive } from 'vue'

export default function useRandomPWD() {
  const password = reactive({
    checkbox: [1, 2, 3],
    length: 8,
    text: '',
  })
  const createPassword = () => {
    let randArr: string[] = []
    for (const num of password.checkbox) {
      if (num === 1) {
        for (let i = 0; i < 26; i++) {
          randArr.push(String.fromCharCode(65 + i))
        }
      } else if (num === 2) {
        for (let i = 0; i < 26; i++) {
          randArr.push(String.fromCharCode(97 + i))
        }
      } else if (num === 3) {
        for (let i = 0; i < 10; i++) {
          randArr.push(i)
        }
      } else {
        randArr = randArr.concat(['!', '@', '#', '$', '%', '^', '&', '*'])
      }
    }
    let tmpPWD = ''
    for (let i = 0; i < password.length; i++) {
      tmpPWD += randArr[Math.floor(Math.random() * randArr.length)]
    }
    password.text = tmpPWD
  }

  return {
    password,
    createPassword,
  }
}
