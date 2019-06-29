import React from 'react';
import {fetchTasks} from "../redux/actions/tasks";
import {setVisible} from "../redux/actions/tasks";
import {connect} from "react-redux";
import Tasks from "../components/Tasks";
import Grid from "@material-ui/core/Grid/Grid";
import Snackbar from "@material-ui/core/Snackbar/Snackbar";

class TasksContainer extends React.Component {

    state = {
        open: false
    };

    constructor(props) {
        super(props);
        this.props.fetchTasks();


    }

    onButtonClick = (task) => {

        this.setState({
            open: true
        });

        this.props.setVisible(task);
        this.props.fetchTasks();
    };

    handleCloseNotification = () => {
        this.setState({
            open: false
        });
    };

    render() {


        return (
            <React.Fragment>
                <Grid container spacing={24} alignContent={'space-around'}>
                    {this.props.tasks.map((task, i) =>
                        <Grid item xs={6} md={4} lg={3} key={i}>
                            {<Tasks task={task} onClick={() => this.onButtonClick(task)}/>}
                        </Grid>
                    )}
                </Grid>
                <Snackbar
                    anchorOrigin={{
                        vertical: 'bottom',
                        horizontal: 'right',
                    }}
                    open={this.state.open}
                    autoHideDuration={1000}
                    onClose={this.handleCloseNotification}
                    message={<span>This tasks is visible</span>}
                />
            </React.Fragment>
        )
    }
}

const mapStateToProps = (states) => {


    const {tasks} = states.tasks;

    console.log(tasks);

    return {
        tasks
    };
};

const mapsDispatchToProps = (dispatch) => {
    return {
        fetchTasks: () => dispatch(fetchTasks(dispatch)),
        setVisible: (data) => dispatch(setVisible(data, dispatch))
    }
};

export default connect(mapStateToProps, mapsDispatchToProps)(TasksContainer);