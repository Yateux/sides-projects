export const receiveTasks = (data) => {
    return {
        type: 'TASKS',
        payload: data
    }
};

export const fetchTasks = (dispatch) => {


    fetch("http://localhost:8080/api/v1/tasks",
        {
            method: 'GET',
            mode: "cors",
            headers: {
                'Content-Type': 'application/json',
            }
        })
        .then(response => response.json())
        .then(data => dispatch(receiveTasks(data)))
        .catch(error => console.log(error));


    return {
        type: 'RECEIVE_TASKS',
        payload: {}
    }
};


export const setVisible = (data, dispatch) => {

    let visible = false;

    if(data.visible === false) {
        visible = true;
    }
    data = {
        id: data.id,
        visible: visible,
    };

    fetch('http://localhost:8080/api/v1/tasks',
        {
            method: 'PUT',
            mode: 'cors',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
        })
        .then(response => response.json())
        .then(data => console.log("ok"))
        .catch(error => console.log(error));

    return {
        type: 'RECEIVE_TASKS',
        payload: {}
    }
};