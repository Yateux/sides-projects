const reducer = function (state = { tasks: [] }, action) {

    switch (action.type) {
        case 'TASKS':
            return Object.assign({}, state, {
                tasks: action.payload
            });

        default:
            return state;
    }
};

export default reducer;