import React from 'react';
import Typography from "@material-ui/core/Typography/Typography";
import Grid from "@material-ui/core/Grid";
import Tasks from "../components/Tasks";
import {fetchTasks, setVisible} from "../redux/actions/tasks";
import {connect} from "react-redux";


 class Home extends React.Component {


     constructor(props) {
         super(props);
         this.props.fetchTasks();


     }

    render() {
        return <div>
            <Typography align="center" variant="h2">
                Tasks
            </Typography>
            <hr/>
            <div className='content'>


                <React.Fragment>
                    <Grid container spacing={24} alignContent={'space-around'}>
                        {this.props.tasks.map((task, i) =>
                            <Grid item xs={6} md={4} lg={3} key={i}>
                                {!task.visible && <Tasks home={true} task={task}/>}
                            </Grid>
                        )}
                    </Grid>
                </React.Fragment>

            </div>
        </div>
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

export default connect(mapStateToProps, mapsDispatchToProps)(Home);