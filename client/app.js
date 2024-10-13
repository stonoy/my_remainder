import React from "react"
import {useEffect, useState} from "react"
import ReactDOM from "react-dom/client"
import {Provider} from "react-redux"
import store from "./store"
import { useDispatch, useSelector } from "react-redux"
import {  loginAsync, setDefaultConfig } from "./features/user/userSlice"
import { createBrowserRouter, Link, Navigate, RouterProvider } from "react-router-dom"
import { Outlet } from "react-router-dom"
import { useNavigate } from "react-router-dom"
import { getRemainders } from "./features/remainders/remainderSlice"
import { deleteRemainder, getTheRemainder, startTimer, updateRemainder } from "./features/remainder/remainderSlice"

const heading = React.createElement("h1", {}, "Welcome!")
               
const Login = () => {
    const {submitting,success} = useSelector(state => state.user)
    const dispatch = useDispatch()
    const formRef = React.useRef(null)
    const navigate = useNavigate()

    useEffect(() => {
        if (success) {
            navigate("/")
        }

        return () => {
            dispatch(setDefaultConfig())
        }
    },[success])

    const handleLogin = (e) => {
        e.preventDefault()

        formData = new FormData(formRef.current)
        data = Object.fromEntries(formData)

        console.log(data)
        dispatch(loginAsync(data))
        
    }

    return (
        <div>
            <form ref={formRef} onSubmit={handleLogin}>
            <div>
                <label>Email</label>
                <input type="text" name="email"/>
            </div>
            <div>
                <label>Password</label>
                <input type="text" name="password"/>
            </div>
            <button type="submit" disabled={submitting}>
                {submitting ? "submitting" : "submit"}
            </button>
            </form>
            <p>New member <Link to="/register">Register</Link></p>
        </div>
    )
}

const Register = () => {


    return (
        <div>
            <div>
                <label>Username</label>
                <input type="text" name="name"/>
            </div>
            <div>
                <label>Email</label>
                <input type="text" name="email"/>
            </div>
            <div>
                <label>Password</label>
                <input type="text" name="password"/>
            </div>
            <button>Submit</button>
            <p>Alreacy a member <Link to="/login">Login</Link></p>
        </div>
    )
}

const Navbar = () => {
    const {token, name} = useSelector(state => state.user)

    return (
        <div>
            <h2>Remainder</h2>
            <ul>
                <li>Contact Us</li>
                <li>{token ? `${name}` : "Login/Register"}</li>
            </ul>
        </div>
    )
}

const HomeLayout = () => {


    return (
        <>
            <Navbar/>
            <hr/>
            <Outlet/>
            
        </>
    )
}

const Contact = () => {
    return (
        <h1>Hi, we are from Remainders...</h1>
    )
}

const Remainders = () => {
    const dispatch = useDispatch()
    const {remainders, loading} = useSelector(state => state.remainder)
    const [showTimeID, setShowTimeID] = useState("")
    const [stopTimer, setStopTimer] = useState(true)

    useEffect(() => {
        dispatch(getRemainders())
    }, [])

    if (loading) {
        return <h1>Loading...</h1>
    }

    if (remainders.length == 0){
        return <h1>no remainders to show</h1>
    }

    return (
        <>
        <button onClick={() => setStopTimer(prev => !prev)}>{stopTimer ? "on" : "off"}</button>
            {remainders.map((remainder) => {
                return <Remainder {...remainder} showTimeID={showTimeID} setShowTimeID={setShowTimeID} id={remainder.id}/>
                
            })}
            {!stopTimer && <Timer/>}
        </>
    )
}

const Remainder = ({id,description, has_priority, timing, subject, day, hour, miniute, second, showTimeID, setShowTimeID}) => {
    

    return (
        <div onClick={() => setShowTimeID(id)}>
                    <h1>{subject}</h1>
                <p>{timing}</p>
                <hr/>
                {showTimeID == id && <Rdetails id={id} has_priority={has_priority} subject={subject} timing={timing} description={description}/>}
                </div>
    )
}

const Rdetails = ({id, description, has_priority, subject, timing}) => {
    const {remainder:{day,hour, miniute, second}, loading, updating, deleting, deleted} = useSelector(state => state.theRemainder)
    const dispatch = useDispatch()
    const formRef = React.useRef(null)
    const navigate = useNavigate()

    useEffect(() => {
        dispatch(getTheRemainder(id))
    }, [])

    useEffect(() => {
        if (deleted) {
            
            navigate(0)
        }
    }, [deleted])

    const handleUpdate = (e) => {
        e.preventDefault()

        const formData = new FormData(formRef.current)
        const data = Object.fromEntries(formData)

        dispatch(updateRemainder({...data,id, has_priority: e.target.has_priority.checked}))
    }

    if (loading) {
        return <h2 style={{color: "blue"}}>Loading...</h2>
    }

    return (
        <div>
            <ul>
                    <li>{day}</li>
                    <li>{hour}</li>
                    <li>{miniute}</li>
                    <li>{second}</li>
                </ul>
                <form ref={formRef} onSubmit={handleUpdate}>
                    <label>Subject</label>
                    <input type="text" name="subject" defaultValue={subject}/>
                    <label>Description</label>
                    <input type="text" name="description" defaultValue={description}/>
                    <label>Timing</label>
                    <input type="datetime-local" name="timing" defaultValue={timing.slice(0,-4)}/>
                    <label>Priority</label>
                    <input type="checkbox" name="has_priority" defaultChecked={has_priority}/>
                    <button type="submit" disabled={updating}>{updating ? "updating..." : "update"}</button>
                </form>

                <br/>
                <button onClick={() => dispatch(deleteRemainder(id))} disabled={deleting} style={{color: "red"}}>Delete</button>
        </div>
    )
}

const Timer = () => {
    const dispatch = useDispatch()

    useEffect(() => {
        const id = setInterval(() => {
            console.log("ticks")
            dispatch(startTimer())
        }, 1000)

        return () => {
            clearTimeout(id)
        }
    }, [])
}

const App = function(){
    const router = createBrowserRouter(
        [
            {
                path: "/login",
                element: <Login/>,
            },
            {
                path: "/register",
                element: <Register/>
            },
            {
                path: "/",
                element: <HomeLayout/>,
                children: [
                    {
                        index: true,
                        element: <Remainders/>
                    },
                    {
                        path: "/contact",
                        element: <Contact/>
                    }
                ]
            }
        ]
    )

    return (
        <RouterProvider router={router}/>
    )
}

const root = ReactDOM.createRoot(document.getElementById("root"))
root.render(
    <Provider store={store}>
        <App/>
    </Provider>
)