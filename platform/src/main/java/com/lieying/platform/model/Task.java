package com.lieying.platform.model;

import jakarta.persistence.*;
import lombok.Data;
import org.hibernate.annotations.CreationTimestamp;
import org.hibernate.annotations.UpdateTimestamp;

import java.time.LocalDateTime;
import java.util.UUID;

@Data
@Entity
@Table(name = "tasks")
public class Task {
    @Id
    @GeneratedValue
    private UUID id;

    @Column(nullable = false)
    private String name;

    @Enumerated(EnumType.STRING)
    @Column(nullable = false)
    private TaskType type;

    @Enumerated(EnumType.STRING)
    @Column(nullable = false)
    private TaskStatus status = TaskStatus.PENDING;

    @Column(name = "target_id")
    private UUID targetId;

    @Column(columnDefinition = "jsonb")
    private String config;

    private Integer progress = 0;

    @Column(name = "result_count")
    private Integer resultCount = 0;

    @Column(name = "started_at")
    private LocalDateTime startedAt;

    @Column(name = "completed_at")
    private LocalDateTime completedAt;

    @CreationTimestamp
    @Column(name = "created_at", nullable = false, updatable = false)
    private LocalDateTime createdAt;

    @UpdateTimestamp
    @Column(name = "updated_at")
    private LocalDateTime updatedAt;

    public enum TaskType {
        RECON, SCAN, FUZZ, EXPLOIT, FULL_CHAIN
    }

    public enum TaskStatus {
        PENDING, RUNNING, PAUSED, COMPLETED, FAILED
    }
}
