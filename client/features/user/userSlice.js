import {createSlice} from "@reduxjs/toolkit"
import { createAsyncThunk } from "@reduxjs/toolkit"
import { axiosBase } from "../../src/utils"

const defaultState = {
    submitting: false,
    success: false,
    name : "Guest",
    token: ""
}

export const loginAsync = createAsyncThunk(
    "user/loginAsync",
    async (data, thunkAPI) => {
        try {
            
            const resp = await axiosBase.post("/login", data)
            
            return resp.data
        } catch (error) {
            
            return thunkAPI.rejectWithValue(error)
        }
    }
)

const userSlice = createSlice({
    name: "user",
    initialState : JSON.parse(localStorage.getItem("user")) || defaultState,
    reducers: {
        logOutUser : () => {
            return defaultState
        },
        setDefaultConfig : (state) => {
            state.submitting = false
            state.success = false
            localStorage.setItem("user", JSON.stringify(state))
        }
    },
    extraReducers: (builder) => {
        builder.addCase(loginAsync.pending, (state, {payload}) => {
            state.submitting = true
        }).addCase(loginAsync.fulfilled, (state, {payload}) => {
            state.submitting = false
            state.success = true
            state.token = payload.token
            state.name = payload.user.name
            localStorage.setItem("user", JSON.stringify(state))
        }).addCase(loginAsync.rejected, (state, action) => {
            state.submitting = false
            
            console.log(action)
        })
    }
})

export const { logOutUser, setDefaultConfig} = userSlice.actions

export default userSlice.reducer