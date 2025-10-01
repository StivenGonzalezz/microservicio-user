package org.example.model;

import jdk.jshell.Snippet;
import lombok.Builder;
import lombok.Getter;
import lombok.Setter;
import lombok.ToString;
import org.example.model.enums.Channels;
import org.example.model.enums.Status;
import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.mapping.Document;
import java.time.LocalDateTime;
import java.util.List;

@Document(collection = "notifications")
@Getter
@Setter
@ToString
public class Notification {
    @Id
    private String id;
    private String message;
    private String affair;
    private String channels;
    private String email;
    private String number;
    private Status status; // PENDING, SENT, FAILED
    private LocalDateTime createdAt;
    private LocalDateTime sendAt;

    @Builder
    public Notification (String message,Status status,LocalDateTime createdAt){

        this.message = message;
        this.status = status;
        this.createdAt = createdAt;

    }

}
