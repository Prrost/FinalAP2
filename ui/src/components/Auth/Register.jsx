import React, { useState } from 'react'

export default function Register() {
    const [firstName, setFirstName] = useState()
    const [lastName, setLastName] = useState()    
    const [email, setEmail] = useState()
    const [password, setPassword] = useState()

  return (
    <div>
    <h1>Register</h1>
    <form>
      <input type='text' placeholder='First Name' />
      <input type='text' placeholder='Last Name' />
      <input type='email' placeholder='Email' />
      <input type='password' placeholder='Password' />
    </form>
  </div>
  )
}
