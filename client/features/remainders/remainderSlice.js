import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import { axiosBase, extractDayTime, ticking } from "../../src/utils";

const defaultState = {
    remainders: [],
    loading: false,
}

export const getRemainders = createAsyncThunk(
    "remainder/getRemainders",
    async (_,thunkAPI) => {
        const token = thunkAPI.getState().user.token
        
        try {
            const resp = await axiosBase.get("/getremainders", 
                {
                    headers : {
                        "Authorization" : `Bearer ${token}`
                    }
                }
            )
            return resp?.data
        } catch (error) {
            return thunkAPI.rejectWithValue(error)
        }
        
    }
)

const remainderSlice = createSlice({
    name : "remainder",
    initialState: defaultState,
    reducers: {
        setRemainders : (state, {payload}) => {
            state.remainders = state.remainders.map(remainder => {
                if (remainder.id == payload.id){
                    return payload
                } else {
                    return remainder
                }
            })
        },
        
    },
    extraReducers : (builder) => {
        builder.addCase(getRemainders.pending, (state) => {
            state.loading = true
        }).addCase(getRemainders.fulfilled, (state, {payload})=> {
            state.loading = false
            state.remainders = payload.remainders.map((remainder) => {
                const dayTime = extractDayTime(remainder.timing)
                return {...remainder, ...dayTime}
            })
            // console.log(payload)
        }).addCase(getRemainders.rejected, (state, action) => {
            state.loading = false
            console.log(action)
        })
    }
})

export const {setRemainders} = remainderSlice.actions

export default remainderSlice.reducer