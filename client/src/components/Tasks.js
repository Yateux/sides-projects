import React from "react";
import Card from "@material-ui/core/Card/Card";
import CardHeader from "@material-ui/core/CardHeader/CardHeader";
import CardContent from "@material-ui/core/CardContent/CardContent";
import Typography from "@material-ui/core/Typography/Typography";
import {withStyles} from "@material-ui/core";
import CardActions from "@material-ui/core/CardActions/CardActions";
import Button from "@material-ui/core/Button/Button";
import Grid from "@material-ui/core/Grid";


const styles = {
    media: {
        height: 0,
        paddingTop: '56.25%'
    },
    pos: {
        marginBottom: 12,
    },
};

class Tasks extends React.PureComponent {


    render() {
        const {task, classes} = this.props;

        return (

            <Card>
                <CardHeader title={task.name}/>
                <CardContent>
                    <Typography className={classes.pos} color="textSecondary">
                        {task.Type.category}

                    </Typography>
                    <Typography paragraph>
                        test
                    </Typography>
                </CardContent>
                {this.props.home !== true && <CardActions>
                    <Button size="small" color="primary"
                            onClick={() => this.props.onClick()}> {task.visible === false ? "Set invisible" : "Set visible"}</Button>
                </CardActions>}

            </Card>

        )
    }
}

export default withStyles(styles)(Tasks);
