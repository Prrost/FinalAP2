import React, { useState } from 'react'

export default function Login() {
    const [login, setLogin] = useState()
    const [password, setPassword] = useState()

    const handleLoginButton = () => {}

  return (
    <div>
      <h1>Login</h1>
      <form>
        <input type='text' placeholder='Login' onChange={(e)=> {setLogin(e.target.value)}}/>
        <input type='password' placeholder='Password' onChange={(e) => {setPassword(e.target.value)}}/>
        <button onClick={handleLoginButton}>Login</button>
        <p>Don't have an accout?</p>
        <a>Register</a>
      </form>
    </div>
  )
}
