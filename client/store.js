import {configureStore} from "@reduxjs/toolkit"
import userSliceReducer from './features/user/userSlice'
import remaindersSlice from "./features/remainders/remainderSlice"
import remainderSlice from "./features/remainder/remainderSlice"

export default store = configureStore({
    reducer : {
       user : userSliceReducer,
       remainder : remaindersSlice,
       theRemainder : remainderSlice,
    }
})