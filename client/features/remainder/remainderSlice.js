import {createSlice, createAsyncThunk} from "@reduxjs/toolkit"
import { axiosBase, extractDayTime, ticking } from "../../src/utils"
import { setRemainders } from "../remainders/remainderSlice"

const defaultState = {
    loading: false,
    updating: false,
    deleting: false,
    deleted: false,
    remainder: {
        day: 0,
        hour: 0,
        miniute: 0,
        second: 0
    }
}

export const getTheRemainder = createAsyncThunk(
    "remainder/getTheRemainder",
    async (id, thunkAPI) => {
        const token = thunkAPI.getState().user.token
        try {
            const resp = await axiosBase.get(`/getremainder/${id}`, {
                headers : {
                    "Authorization": `Bearer ${token}`
                }
            })
            return resp?.data
        } catch (error) {
            return thunkAPI.rejectWithValue(error)
        }
    }
)

export const updateRemainder = createAsyncThunk(
    "remainder/updateRemainder",
    async (data, thunkAPI) => {
        const token = thunkAPI.getState().user.token
        try {
            const resp = await axiosBase.put(`/updateremainder/${data.id}`, data, {
                headers : {
                    "Authorization": `Bearer ${token}`
                }
            })
            thunkAPI.dispatch(setRemainders(resp?.data?.remainder))
            return resp?.data
        } catch (error) {
            return thunkAPI.rejectWithValue(error)
        }
    }
)

export const deleteRemainder = createAsyncThunk(
    "remainder/deleteRemainder",
    async (id, thunkAPI) => {
        const token = thunkAPI.getState().user.token
        try {
            const resp = await axiosBase.delete(`/deleteremainder/${id}`, {
                headers : {
                    "Authorization": `Bearer ${token}`
                }
            })
            return resp?.data
        } catch (error) {
            return thunkAPI.rejectWithValue(error)
        }
    }
)

const remainderSlice = createSlice({
    name: "remainder",
    initialState: defaultState,
    reducers: {
        startTimer : (state) => {
            state.remainder = ticking(state.remainder)
        }
    },
    extraReducers : (builder) => {
        builder.addCase(getTheRemainder.pending, (state) => {
            state.loading = true
            state.deleted = false
        }).addCase(getTheRemainder.fulfilled, (state, {payload}) => {
            state.loading = false
            // console.log(payload)
            state.remainder = extractDayTime(payload.remainder.timing)
        }).addCase(getTheRemainder.rejected, (state, action) => {
            state.loading = false
            console.log(action)
        }).addCase(updateRemainder.pending, (state) => {
            state.updating = true
        }).addCase(updateRemainder.fulfilled, (state, {payload}) => {
            state.updating = false
            // console.log(payload)
            state.remainder = extractDayTime(payload.remainder.timing)
        }).addCase(updateRemainder.rejected, (state, action) => {
            state.updating = false
            console.log(action)
        }).addCase(deleteRemainder.pending, (state) => {
            state.deleting = true
        }).addCase(deleteRemainder.fulfilled, (state, {payload}) => {
            state.deleting = false
            state.deleted = true
            // console.log(payload)
        }).addCase(deleteRemainder.rejected, (state, action) => {
            state.deleting = false
            console.log(action)
        })
    }
})

export const {startTimer} = remainderSlice.actions

export default remainderSlice.reducer